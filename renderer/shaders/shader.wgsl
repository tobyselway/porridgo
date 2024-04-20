struct VertexIn {
    @location(0) position: vec3<f32>,
    @location(1) texture: vec2<f32>,
}

struct VertexOutput {
    @builtin(position) position: vec4<f32>,
    @location(1) texture : vec2<f32>,
}

@vertex
fn vs_main(in: VertexIn) -> VertexOutput {
    var out: VertexOutput;
    // TODO: Transform from world space to clip space
    out.position = vec4<f32>(in.position, 1.0);
    out.texture = in.texture;
    return out;
}

@group(0)
@binding(0)
var image: texture_2d<f32>;

@fragment
fn fs_main(vertex : VertexOutput) -> @location(0) vec4<f32> {
    return textureLoad(image, vec2<i32>(vertex.texture * 256.0), 0);
}
