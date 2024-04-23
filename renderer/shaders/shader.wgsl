struct VertexInput {
    @location(0) position: vec3<f32>,
    @location(1) uv_mapping: vec2<f32>,
}

struct VertexOutput {
    @builtin(position) clip_position: vec4<f32>,
    @location(1) uv_mapping : vec2<f32>,
}

@vertex
fn vs_main(in: VertexInput) -> VertexOutput {
    var out: VertexOutput;
    out.uv_mapping = in.uv_mapping;
    out.clip_position = vec4<f32>(in.position, 1.0);
    return out;
}

@group(0)
@binding(0)
var t_diffuse: texture_2d<f32>;

@group(0)
@binding(1)
var s_diffuse: sampler;

@fragment
fn fs_main(vertex : VertexOutput) -> @location(0) vec4<f32> {
    return textureSample(t_diffuse, s_diffuse, vertex.uv_mapping);
}
